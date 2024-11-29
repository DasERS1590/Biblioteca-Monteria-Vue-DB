import React from "react";
import { Navigate } from "react-router-dom";

const DashboardRedirect = ({ userRole }) => {
  return (
    <div>
      {userRole === "admin" ? <Navigate to="/admin/dashboard" /> : <Navigate to="/user/dashboard" />}
    </div>
  );
};

export default DashboardRedirect;
