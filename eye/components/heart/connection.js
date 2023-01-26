class Connection {
  constructor() {
    this.callbacks = {
      onconnect: [],
      onmessage: [],
      onclose: [],
      onerror: [],
    }
  }
  connect(url) {
    const self = this
    if (process.server === true) {
      return
    }
    console.log("❤️: Starting connection to WebSocket Server on ", url)
    if (this.connection) {
      this.connection.close()
    }
    this.connection = new WebSocket(url)

    this.connection.onerror = function (error) {
      self.emit('onerror', error)
    }

    this.connection.onended = function (e) {
      self.emit('onended', null)
    }

    this.connection.onabort = function (e) {
      self.emit('onclose', null)
    }

    this.connection.onclose = function (e) {
      self.emit('onclose', null)
    }

    this.connection.onmessage = function (event) {
      event = JSON.parse(event.data);
      if (event.Name === 'creation') { // card is in work now
        let creationState = JSON.parse(event.Payload)
        self.emit('onmessage', creationState)
      }
    }

    this.connection.onopen = function (event) {
      console.log(`Successfully connected to the echo websocket server ${url}`)
      self.emit('onconnect', null)
    }
  }
  close() {
    console.log('close')
    this.connection.close()
  }
  onconnect(cb) {
    this.callbacks.onconnect.push(cb)
  }
  onmessage(cb) {
    this.callbacks.onmessage.push(cb)
  }
  onclose(cb) {
    this.callbacks.onclose.push(cb)
  }
  onerror(cb) {
    this.callbacks.onerror.push(cb)
  }
  emit(type, event) {
    let callbacks = this.callbacks[type]
    if (callbacks.length === 0) {
      return
    }
    callbacks.forEach((cb) => {
      cb(event)
    })
  }
}
export default Connection
