import React from 'react';
import './App.css';

function App() {
  let socket = new WebSocket("ws://" + document.location.hostname + (document.location.port ? ':' + document.location.port : '') + "/ws");

  socket.onopen = (() => {
    console.log("connection")
  })

  socket.onmessage = ((message) => {
    console.log(message)
  })

  return (
    <div className="App">
      <header className="App-header">
        Hello
      </header>
    </div>
  );
}

export default App;
