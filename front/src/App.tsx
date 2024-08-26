import React from 'react';
import './App.css';

function App() {
  let ws = new WebSocket("ws://" + document.location.hostname + (document.location.port ? ':' + document.location.port : '') + "/ws");

  ws.onopen = () => {
    console.log("connection success");
    ws.send("Hello from browser!");
  };
  
  ws.onmessage = (event) => {
    console.log("message:", event.data);
  }

  return (
    <div className="App">
      <header className="App-header">
        Hello
      </header>
    </div>
  );
}

export default App;
