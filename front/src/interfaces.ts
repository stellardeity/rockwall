export type DialogType = {
  date: string;
  isMine: boolean;
  from: string;
  content: string;
};

export interface IState {
  iam: {
    id: string;
    name: string;
  } | null;
  interlocutor: {
    id: string;
    name: string;
  } | null;
  messages: {
    [key: string]: DialogType[];
  };
}

export type PeerType = {
  counter: number;
  id: string;
  name: string;
};

export type PeersType = {
  [key: string]: PeerType;
};

export interface ISocketMessage {
  name: string;
  id: string;
  cmd: "NAME" | "MESS" | "PEERS";
  from: string;
  to: string;
  peers: PeerType[];
  content: string;
};
