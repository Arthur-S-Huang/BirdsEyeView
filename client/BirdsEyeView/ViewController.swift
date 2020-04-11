//
//  ViewController.swift
//  BirdsEyeView
//
//  Created by Jade Lundy on 12/11/19.
//  Copyright © 2019 BEV196. All rights reserved.
//

import UIKit

class ViewController: UIViewController, UINavigationControllerDelegate, UIImagePickerControllerDelegate {

    override func viewDidLoad() {
        super.viewDidLoad()
        // Do any additional setup after loading the view.
    }
    
    var imagePicker: UIImagePickerController!
    
    @IBOutlet weak var imageView: UIImageView!
    
    @IBAction func takePhoto(_ sender: Any) {
        
        if !UIImagePickerController.isSourceTypeAvailable(.camera){

            let alertController = UIAlertController.init(title: nil, message: "Device has no camera.", preferredStyle: .alert)

            let okAction = UIAlertAction.init(title: "Alright", style: .default, handler: {(alert: UIAlertAction!) in
            })

            alertController.addAction(okAction)
            self.present(alertController, animated: true, completion: nil)

        }
        else{
             imagePicker =  UIImagePickerController()
             imagePicker.delegate = self
             imagePicker.sourceType = .camera

             present(imagePicker, animated: true, completion: nil)
        }
    }
    
    
    
    
    func imagePickerController(_ picker: UIImagePickerController, didFinishPickingMediaWithInfo info: [UIImagePickerController.InfoKey : Any]) {
        imagePicker.dismiss(animated: true, completion: nil)
        
        let image =  info[.originalImage] as? UIImage
        imageView.image = image
        
        
        let completionHandler: (Bool, Error?) -> Void = { _,_ in } // TODO Handle completion
        //Make API call to get presigned url
        upload(image: image!, urlString: "https://birdseyeview.s3.amazonaws.com/test2.jpg?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=AKIA4OV4J5XFFV5ZRQBQ%2F20191212%2Fus-east-1%2Fs3%2Faws4_request&X-Amz-Date=20191212T052037Z&X-Amz-Expires=18000&X-Amz-SignedHeaders=host&X-Amz-Signature=00b701cac30cc6ffaf6efb05470cc1d00223994ea58e823fac54d121dd19fdaf", mimeType: "image/jpeg", completion: completionHandler)
        
        let url = NSURL(string: "https://arthur.local:19603/caption/test2.jpg")
        let request = NSMutableURLRequest(url: url! as URL)
        request.httpMethod = "GET"
        
        let session = URLSession(configuration:URLSessionConfiguration.default, delegate: nil, delegateQueue: nil)
        let dataTask = session.dataTask(with: request as URLRequest)
        dataTask.resume()
    }
    
    func upload(data: Data, urlString: String, mimeType: String, completion: @escaping (Bool, Error?) -> Void) {
        let url = NSURL(string: urlString)
        let request = NSMutableURLRequest(url: url! as URL)
        request.httpMethod = "PUT"
        request.httpBody = data
        
        let session = URLSession(configuration:URLSessionConfiguration.default, delegate: nil, delegateQueue: nil)
        let dataTask = session.dataTask(with: request as URLRequest)
        dataTask.resume()
        
        
    }
        
    func upload(image: UIImage, urlString: String, mimeType: String, completion: @escaping (Bool, Error?) -> Void) {
        let data = image.jpegData(compressionQuality: 0.9)!
        upload(data: data, urlString: urlString, mimeType: mimeType, completion: completion)
    }

    
}
