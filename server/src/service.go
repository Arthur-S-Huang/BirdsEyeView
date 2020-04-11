package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/rekognition"
    "github.com/aws/aws-sdk-go/service/polly"
    "github.com/aws/aws-sdk-go/service/sqs"
    //"github.com/aws/aws-sdk-go/service/mediaconvert"
    "github.com/aws/aws-sdk-go/service/s3"
    "github.com/buger/jsonparser"
    "strings"
    "net/http"
    "encoding/json"
    "io/ioutil"
    "fmt"
    "time"
    "github.com/gorilla/mux"
)

//need to implement for video
var mybucket string = "birdseyeview"
var myphoto string = ""
var myvideo string = ""
var key string = ""
var welcome string = "Hello there! Let me tell you some of the things I can recognize in your image. I can see: "
//var id = result
//var videoResults string = ""
func configureBucket(bucketName string, photoName string) {
    mybucket = bucketName
    myphoto = photoName
}

func createImageCaption() []string {
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    svc := rekognition.New(sess)
    input := &rekognition.DetectLabelsInput{
        Image: &rekognition.Image{
            S3Object: &rekognition.S3Object{
                Bucket: aws.String(mybucket),
                Name:   aws.String(myphoto),
            },
        },
        MaxLabels:     aws.Int64(20),
        MinConfidence: aws.Float64(80.000000),
    }
    result, _ := svc.DetectLabels(input)

    ot := make([]string, 0)

    for _, label := range result.Labels {
        ot = append(ot, *(label.Name))
    }
    return ot
}

func createImageSpeech() []byte{
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))

    // Create Polly client
    svc := polly.New(sess)
    b := createImageCaption()
    s := strings.Join(b, ", ")
    ss := welcome + s
    //fmt.Println(s)
    // Output to MP3 using voice Joanna
    input := &polly.SynthesizeSpeechInput{Engine: aws.String("neural"),
                                          OutputFormat: aws.String("mp3"),
                                          Text: aws.String(ss),
                                          VoiceId: aws.String("Salli")}
    output, _ := svc.SynthesizeSpeech(input)

    //fmt.Println(output)

    a := output.AudioStream
    body, _ := ioutil.ReadAll(a)

    return body
}
//check for correct JobId
//use MediaConvert to convert input to mp4, rotate input 90 degrees clockwise
func createVideoCaption(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    myvideo = params["name"]
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
/*
    sess2 := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    */
    svc := rekognition.New(sess)

    //svc2 := polly.New(sess2)
    //convertVideo()
    //time.Sleep(90 * time.Second)
    input := &rekognition.StartLabelDetectionInput{
        NotificationChannel: &rekognition.NotificationChannel{
            RoleArn: aws.String("arn:aws:iam::856167149002:role/videoAnalysis"),
            SNSTopicArn: aws.String("arn:aws:sns:us-east-1:856167149002:birdseyeview"),
        },
        Video: &rekognition.Video{
            S3Object: &rekognition.S3Object{
                Bucket: aws.String(mybucket),
                Name:   aws.String(myvideo),
            },
        },
        MinConfidence: aws.Float64(80.000000),
    }
    fmt.Println(input)
    result, _ := svc.StartLabelDetection(input)

    //getLabelInput := &rekognition.GetLabelDetectionInput{JobId: aws.String(*id.JobId)}
    fmt.Println(result)
    ot := make([]string, 0)
    for {
        time.Sleep(5 * time.Second)
        complete := checkDetection(*result.JobId)
        if (complete == true) {
            fmt.Println("complete")
            labels, _ := svc.GetLabelDetection(&rekognition.GetLabelDetectionInput{JobId: aws.String(*result.JobId),
                                                                                   MaxResults: aws.Int64(10)})
            for _, label := range labels.Labels {
                if (contains(ot, *label.Label.Name)) {
                    continue
                } else {
                    ot = append(ot, *(label.Label.Name))
                }
            }
            w.Write(tts(ot))
            //break
            //fmt.Println(ot)
            //json.NewEncoder(w).Encode(ot)
            break
        }
    }
}

func checkDetection(id string) bool {
    var completed bool = false
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    svc := sqs.New(sess)

    messages, _ := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
        QueueUrl:            aws.String("https://sqs.us-east-1.amazonaws.com/856167149002/birdseyeview"),
        MaxNumberOfMessages: aws.Int64(10),
        VisibilityTimeout:   aws.Int64(0),
        WaitTimeSeconds:     aws.Int64(0),
    })

    for _, message := range messages.Messages {
        notification := []byte(*message.Body)
        //fmt.Println(*message.Body)
        msg, _, _, _ := jsonparser.Get(notification, "Message")
        //fmt.Println(string(msg))

        if (strings.Contains(string(msg), "SUCCEEDED") && strings.Contains(string(msg), id)) {
            //fmt.Println(*message.Body)
            completed = true
            break
        }

    }
    return completed
}

func tts(txt []string) []byte{
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    videoResults := strings.Join(txt, ", ")
    svc := polly.New(sess)
    read := welcome + videoResults + "Wellll, wasn't that interesting? Let me know what else you want me to see!"
    input := &polly.SynthesizeSpeechInput{Engine: aws.String("neural"),
                                           OutputFormat: aws.String("mp3"),
                                           Text: aws.String(read),
                                           VoiceId: aws.String("Salli")}
    output, _ := svc.SynthesizeSpeech(input)
    a := output.AudioStream
    body, _ := ioutil.ReadAll(a)
    return body
}

func getPresignedUrl(w http.ResponseWriter, r *http.Request) {
    // Initialize a session in us-west-2 that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials.
    params := mux.Vars(r)
    key = params["file"]
    sess := session.Must(session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    }))
    // Create S3 service client
    svc := s3.New(sess)

    req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
        Bucket: aws.String("birdseyeview"),
        Key:    aws.String(key),
    })
    str, _ := req.Presign(300 * time.Minute)

    fmt.Println("Unencoded url: ", str)
    json.NewEncoder(w).Encode(str)
}

func contains(arr []string, str string) bool {
    for _, a := range arr {
        if a == str {
            return true
        }
    }
    return false
}

func getImageCaption(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    myphoto = params["name"]
    fmt.Println(myphoto)
    b := createImageCaption()
    json.NewEncoder(w).Encode(b)
}

func getImageSpeech(w http.ResponseWriter, r *http.Request) {
    params := mux.Vars(r)
    myphoto = params["name"]
    b := createImageSpeech()
    w.Write(b)
}
