import React from "react";
import Navbar from "../Navbar";
import "../../styles/User.css";

const UserLayout = ({ children }) => {
  return (
    <div className="user-layout">
      <Navbar isAdmin={false} />
      <div className="content">{children}</div>
    </div>
  );
};

export default UserLayout;
