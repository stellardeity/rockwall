import React, { FC } from "react";

interface IPeerProps {
  id: string;
  peers: any;
  selectPeer: (peer: any) => void;
}

export const Peer: FC<IPeerProps> = ({ id, peers, selectPeer }) => {
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
};
