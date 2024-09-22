import React, { FC } from "react";

interface IMessageProps {
  date: string;
  from: string;
  content: string;
}

export const Message: FC<IMessageProps> = ({ date, from, content }) => {
  return (
    <div key={date}>
      <p>{from}</p>
      <p>{content}</p>
    </div>
  );
};
