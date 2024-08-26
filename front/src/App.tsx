import React, { useState } from 'react';
import { Button, ConfigProvider, Input, Space, theme } from 'antd';

function App() {
  const [isConn, setConn] = useState(false);
  const [input, setInput] = useState("");

  let ws = new WebSocket("ws://" + document.location.hostname + (document.location.port ? ':' + document.location.port : '') + "/ws");

  ws.onopen = () => {
    console.log("connection success");
    setConn(true)
  };
  
  ws.onmessage = (event) => {
    console.log("message:", event.data);
  }

  const handleClick = () => {
    if (isConn) {
      ws.send(JSON.stringify({ "ROCK": input}))
    }
  }
  
  return <ConfigProvider
      theme={{
        algorithm: theme.darkAlgorithm,
      }}
    >
      <Space>
        <Input placeholder="Please Input" value={input} onChange={(e) => setInput(e.target.value)} />
        <Button type="primary" onClick={handleClick}>Submit</Button>
      </Space>
    </ConfigProvider>
}

  export default App;
