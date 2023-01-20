import io
import os
import re
import subprocess

from flask import Flask, render_template, Response, abort, request
from io import BytesIO
from PIL import Image

app = Flask(__name__)


@app.route('/painting', methods=['POST'])
def painting():
    print(request.form['tags'])
    print(request.form['seed'])

    filename = getPaintingFromInvokeAIFilename()

    im = Image.open(filename)
    print(im.format, im.size, im.mode)

    img_byte_arr = io.BytesIO()
    im.save(img_byte_arr, format="JPEG")

    return Response(img_byte_arr.getvalue(), content_type="image/jpeg")


def getPaintingFromInvokeAIFilename():
    prepareFileForInvokeAI()

    pattern = re.compile(".*(\/home\/artchitector\/invokeai\/outputs\/[0-9\.]+png).*")
    filename = None

    # print("###START running invoke.py")
    ret = os.popen(
        '/home/artchitector/invokeai/.venv/bin/python /home/artchitector/invokeai/.venv/bin/invoke.py --from_file "/home/artchitector/invokeai/list.txt"'
        )
    # print("### GOT RESPONSE")
    lines = ret.readlines()
    for line in lines:
        match = pattern.match(line)
        if match is not None:
            filename = match.groups()[0]
            print(f"Found filename: {filename}")

    if filename is not None:
        return filename
    else:
        raise Exception("filename not found")


def prepareFileForInvokeAI():
    tags = request.form['tags']
    seed = request.form['seed']
    lines = []
    lines.append(f'{tags} -S{seed} -W512 -H768 -s50 -U2')
    with open("/home/artchitector/invokeai/list.txt", "w") as text_file:
        for line in lines:
            text_file.write(line)
    text_file.close()


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8083, debug=True)
