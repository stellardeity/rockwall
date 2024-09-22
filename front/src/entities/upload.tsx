/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { FC } from "react";
import { UploadOutlined } from "@ant-design/icons";
import { Button, Upload, UploadProps, UploadFile } from "antd";

export const UploadMusic: FC<{ setFile: (file: UploadFile) => void }> = ({
  setFile,
}) => {
  const onChange: UploadProps["onChange"] = ({ file }) => {
    setFile(file);
  };

  return (
    <Upload name="file" beforeUpload={() => false} onChange={onChange}>
      <Button icon={<UploadOutlined />}>Click to Upload</Button>
    </Upload>
  );
};
