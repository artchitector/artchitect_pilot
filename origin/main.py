import cv2
from flask import Flask, render_template, Response, abort

app = Flask(__name__)
# get system webcam
camera = cv2.VideoCapture(0)

def gen_frames():
    while True:
        success, frame = camera.read()  # read the camera frame
        if not success:
            break
        else:
            ret, buffer = cv2.imencode('.jpg', frame)
            frame = buffer.tobytes()
            yield (b'--frame\r\n'
                   b'Content-Type: image/jpeg\r\n\r\n' + frame + b'\r\n')  # concat frame one by one and show result


@app.route('/stream')
def stream():
    """ Video streaming route. """
    return Response(gen_frames(), mimetype='multipart/x-mixed-replace; boundary=frame')

@app.route('/')
def index():
    success, frame = camera.read()
    if not success:
        abort(404)
    else:
        ret, buffer = cv2.imencode('.jpg', frame)
        frame = buffer.tobytes()
        return Response(frame, content_type="image/jpeg")


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8081, debug=True)
