// MobileNet performs target detection for drones using the 
// SSD MobileNet V1 NN.
// For more info on this model see:
// https://github.com/tensorflow/models/tree/master/research/object_detection
type MobileNet {
    // various attributes...
}

// DetectTargets captures an image of the drone's field of view and feeds
// it to a neural network to detect and classify interesting nearby 
// targets.
func (mn *MobileNet) DetectTargets(d *drone.Drone) ([]*Target, error){
    //...
}