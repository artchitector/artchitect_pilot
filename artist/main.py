import io

from flask import Flask, render_template, Response, abort, request
from io import BytesIO
from PIL import Image

app = Flask(__name__)

@app.route('/painting', methods=['POST'])
def painting():
    print(request.form['idea'])
    print(request.form['tags'])
    print(request.form['seed'])

    # with open('allah.jpg', mode='rb') as file:
    #     fileContent = file.read()
    #     return Response(fileContent, content_type="image/jpeg")

    im = Image.open("allah.jpg")
    print(im.format, im.size, im.mode)
    out = im.resize((512, 512), resample=True)

    img_byte_arr = io.BytesIO()
    out.save(img_byte_arr, format="JPEG")

    return Response(img_byte_arr.getvalue(), content_type="image/jpeg")

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8083, debug=True)