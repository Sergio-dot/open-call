# open-call
OpenCall is a project born for my Software Engineering subject at University.
The general idea is to build a web application that allows teleconferencing without the presence of a server that manages the incoming/outgoing data flow.
This is thanks to WebRTC technology which allows participants to act as peers and therefore to exchange information directly.

## Features
OpenCall has some basic features:
- Room creation & generation of a random UUID for room identification
- Join room through invitation received by the room owner
- Audio stream (Microphone can be toggled ON/OFF)
- Video stream (Camera can be toggled ON/OFF)
- Screen sharing (Can be toggled ON/OFF)
- Real-time chat

## Future work
At the actual state this project is much rudimental. There are tons of features to add taking inspiration from similar web & desktop applications.
Some examples of features could be:
- File sharing
- Emoji integration
- Room persistence

In addition there are some improvements to be made to the backend, e.g. improving peer disconnection/reconnection handling

## Technologies
Built in **Go v1.19** with:
- [Fiber](https://github.com/gofiber/fiber) framework
- [GORM](https://github.com/go-gorm/gorm) ORM library
- [Pion](https://github.com/pion/webrtc) WebRTC implementation
