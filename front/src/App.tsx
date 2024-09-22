import React, { useState } from "react";
import { Button, ConfigProvider, Input, Space, UploadFile } from "antd";
import { Message as Dialog } from "./entities/message";
import { Peer } from "./entities/peer";
import {
  DialogType,
  ISocketMessage,
  IState,
  PeersType,
  PeerType,
} from "./interfaces";
import { UploadMusic } from "./entities/upload";

const socket = new WebSocket(
  "ws://" +
    document.location.hostname +
    (document.location.port ? ":" + document.location.port : "") +
    "/ws"
);

function App() {
  const [msg, setMsg] = useState("");
  const [peers, setPeers] = useState<PeersType>({});
  const [file, setFile] = useState<UploadFile>();
  const [state, setState] = useState<IState>({
    iam: null,
    interlocutor: null,
    messages: {},
  });

  const handler = (msgObj: ISocketMessage) => {
    switch (msgObj.cmd) {
      case "NAME": {
        setState((prev: IState) => ({
          ...prev,
          iam: { name: msgObj.name, id: msgObj.id },
        }));
        break;
      }
      case "PEERS": {
        const peers: PeersType = {};
        msgObj.peers.forEach((p) => {
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

        if (state.iam && msgObj.from === state.iam.id) {
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
          isMine: msgObj.from === state.iam?.id,
          from: fromName,
          content: msgObj.content,
        };

        oldMessages.push(message);

        setState((prev: IState) => ({
          ...prev,
          messages: { [peerId]: oldMessages, ...prev.messages },
        }));
        setPeers((prev: PeersType) => ({
          ...prev,
          [peerId]: {
            counter,
            id: "",
            name: "",
          },
        }));
        break;
      }
      default: {
        console.warn("Unknown cmd: " + msgObj.cmd);
      }
    }
  };

  socket.onopen = () => {
    socket.send(JSON.stringify({ cmd: "HELLO" }));
    socket.send(JSON.stringify({ cmd: "PEERS" }));
  };

  socket.onmessage = (event) => {
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
      from: state.iam?.id,
      to: state.interlocutor?.id,
      content: file?.name || msg,
    });
    socket.send(cmd);
  };

  const selectPeer = (peer: PeerType) => {
    setState((prev: IState) => ({
      ...prev,
      interlocutor: peer,
    }));
    setPeers((prev: PeersType) => ({
      ...prev,
      [peer.id]: { id: "", name: "", counter: 0 },
    }));
  };

  socket.onerror = function (error) {
    console.log("Ошибка " + error);
  };

  socket.onclose = (event) => {
    if (event.wasClean) {
      console.log("Соединение закрыто чuисто");
    } else {
      console.log("Обрыв соединения");
    }
    console.log("Код: " + event.code + " причина: " + event.reason);
  };

  return (
    <ConfigProvider>
      <Space>
        {peers &&
          Object.keys(peers)?.map((id) => {
            return (
              <Peer key={id} id={id} peers={peers} selectPeer={selectPeer} />
            );
          })}
        <p>{state.interlocutor ? " with " + state.interlocutor.name : ""}</p>
        <div>
          {state.messages &&
            Object.values(state.messages)?.map((data: DialogType[]) =>
              data.map((data: DialogType) => (
                <Dialog key={data.date} {...data} />
              ))
            )}
        </div>
        <UploadMusic setFile={setFile} />
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
