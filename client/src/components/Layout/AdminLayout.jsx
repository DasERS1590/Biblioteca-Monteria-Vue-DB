import React from "react";
import Navbar from "../Navbar";
import "../../styles/Admin.css";

const AdminLayout = ({ children }) => {
  return (
    <div className="admin-layout">
      <Navbar isAdmin={true} />
      <div className="content">{children}</div>
    </div>
  );
};

export default AdminLayout;
