import socketio

sio = socketio.Client()


@sio.event
def connect():
    print('Connection established')
    sio.emit("myresponse","ASISH",namespace="/",callback=test)


@sio.event
def msg(data):
    print('message ', data)


@sio.event
def my_message(data):
    print('message received with ', data)
    sio.emit('myresponse', {'response': 'my response'})


@sio.event
def disconnect():
    print('Disconnected from server')

def test():
    print("jey")


sio.connect('http://localhost:9090')


sio.emit("hostCreateNewGame", {"hey": "as"},namespace="/")


sio.wait()


