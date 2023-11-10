// NavigateTo applies any required changes to the drone's speed vector 
// so that its eventual position matches dst.
func (d *Drone) NavigateTo(dst Vec3) error { //... }

// Position returns the current drone position vector.
func (d *Drone) Position() Vec3 { //... }

// Position returns the current drone speed vector.
func (d *Drone) Speed() Vec3 { //... }

// CaptureImage records and returns an image of the drone's field of 
// view using the on-board drone camera.
func (d *Drone) CaptureImage() (*image.RGBA, error) { //... }
