package main
/*
import "github.com/aws/aws-sdk-go/service/rekognition"

type DetectLabelsOutput struct {
	_ struct{} `type:"structure"`

	// Version number of the label detection model that was used to detect labels.
	LabelModelVersion *string `type:"string"`

	// An array of labels for the real-world objects detected.
	Labels []*rekognition.Label `type:"list"`

	// The value of OrientationCorrection is always null.
	//
	// If the input image is in .jpeg format, it might contain exchangeable image
	// file format (Exif) metadata that includes the image's orientation. Amazon
	// Rekognition uses this orientation information to perform image correction.
	// The bounding box coordinates are translated to represent object locations
	// after the orientation information in the Exif metadata is used to correct
	// the image orientation. Images in .png format don't contain Exif metadata.
	//
	// Amazon Rekognition doesn’t perform image correction for images in .png
	// format and .jpeg images without orientation information in the image Exif
	// metadata. The bounding box coordinates aren't translated and represent the
	// object locations before the image is rotated.
	OrientationCorrection *string `type:"string" enum:"OrientationCorrection"`
}
/*
type Label struct {
	_ struct{} `type:"structure"`

	// Level of confidence.
	Confidence *float64 `type:"float"`

	// If Label represents an object, Instances contains the bounding boxes for
	// each instance of the detected object. Bounding boxes are returned for common
	// object labels such as people, cars, furniture, apparel or pets.
	Instances []*Instance `type:"list"`

	// The name (label) of the object or scene.
	Name *string `type:"string"`

	// The parent labels for a label. The response includes all ancestor labels.
	Parents []*Parent `type:"list"`
}
*/
