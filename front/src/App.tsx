/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { useState } from "react";
import { Button, ConfigProvider, Input, Space } from "antd";

interface IState {
  socket: any;
  iam: any;
  interlocutor: any;
  messages: any;
}

type MessageType = {
  name: string;
  id: string;
  cmd: "NAME" | "MESS" | "PEERS";
  from: string;
  to: string;
  peers: any;
  content: string;
};

const socket = new WebSocket(
  "ws://" +
    document.location.hostname +
    (document.location.port ? ":" + document.location.port : "") +
    "/ws"
);

function App() {
  const [msg, setMsg] = useState("");
  const [state, setState] = useState<IState>({
    socket: null,
    iam: null,
    interlocutor: null,
    messages: [],
  });
  const [peers, setPeers] = useState<any>([]);

  const handler = (msgObj: MessageType) => {
    switch (msgObj.cmd) {
      case "NAME": {
        setState((prev: IState) => ({
          ...prev,
          iam: { name: msgObj.name, id: msgObj.id },
        }));
        break;
      }
      case "PEERS": {
        const peers: MessageType["peers"] = {};
        msgObj.peers.forEach((p: MessageType["peers"]) => {
          const v = peers[p.id];
          p.counter = v ? v.counter : 0;
          peers[p.id] = p;
        });
        setPeers(peers);
        return;
      }
      case "MESS": {
        let peerId = "";
        let fromName = "";
        let counter = 0;

        if (msgObj.from === state.iam.id) {
          peerId = msgObj.to;
          fromName = state.iam.name;
        } else {
          peerId = msgObj.from;
          const peer = peers[peerId];
          if (peer) {
            fromName = peer.name;
            counter = peer.counter + 1;
          } else {
            fromName = msgObj.from.substr(0, 10);
          }
        }

        let oldMessages = state.messages[peerId];
        if (!oldMessages) {
          oldMessages = [];
        }

        const message = {
          date: new Date().toLocaleTimeString(["ru-RU", "en-US"], {
            hour12: false,
          }),
          isMine: msgObj.from === state.iam.id,
          from: fromName,
          content: msgObj.content,
        };

        oldMessages.push(message);

        setState((prev: IState) => ({
          ...prev,
          messages: { [peerId]: oldMessages, ...prev.messages },
        }));
        setPeers((prev: any) => ({ ...prev, [peerId]: counter }));
        break;
      }
      default: {
        console.warn("Unknown cmd: " + msgObj.cmd);
      }
    }
  };

  socket.onopen = () => {
    console.log("Соединение установлено.");
    socket.send(JSON.stringify({ cmd: "HELLO" }));
    socket.send(JSON.stringify({ cmd: "PEERS" }));
  };

  socket.onmessage = (event) => {
    console.log("Получены данные " + event.data);
    const parsedMessage = JSON.parse(event.data);

    if (!parsedMessage.cmd) {
      console.error("something wrong with data");
      return;
    }

    handler(parsedMessage);
  };

  const sendMessage = () => {
    const cmd = JSON.stringify({
      cmd: "MESS",
      from: state.iam.id,
      to: state.interlocutor.id,
      content: msg,
    });
    socket.send(cmd);
  };

  const selectPeer = (peer: any) => {
    setState((prev: IState) => ({
      ...prev,
      interlocutor: peer,
    }));
    setPeers((prev: any) => ({ ...prev, [peer.id]: 0 }));
  };

  socket.onerror = function (error: any) {
    console.log("Ошибка " + error.message);
  };

  socket.onclose = (event) => {
    if (event.wasClean) {
      console.log("Соединение закрыто чuисто");
    } else {
      console.log("Обрыв соединения");
    }
    console.log("Код: " + event.code + " причина: " + event.reason);
  };

  const interlocutorName = state.interlocutor
    ? " with " + state.interlocutor.name
    : "";

  return (
    <ConfigProvider>
      <Space>
        {peers &&
          Object.keys(peers)?.map((id) => {
            return (
              <div
                key={id}
                data-name={peers[id].name}
                data-id={id}
                onClick={() => selectPeer({ id, name: peers[id].name })}
              >
                <div>{peers[id].name}</div>
                <div>
                  <div>{peers[id].counter > 0 ? peers[id].counter : ""}</div>
                </div>
              </div>
            );
          })}
        <p>{interlocutorName}</p>
        <div>
          {state.messages &&
            Object.values(state.messages)?.map((data: any) => {
              return data.map((data: any) => (
                <div key={data.date}><p>{data.from}</p><p>{data.content}</p></div>
              ));
            })}
        </div>
        <Input
          placeholder="Please Input"
          value={msg}
          onChange={(e) => setMsg(e.target.value)}
        />
        <Button type="primary" onClick={sendMessage}>
          Submit
        </Button>
      </Space>
    </ConfigProvider>
  );
}

export default App;
