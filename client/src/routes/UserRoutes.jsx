import React from "react";
import { Route, Routes } from "react-router-dom";

import UserLayout from "../components/Layout/UserLayout";
import ProtectedRoute from "./ProtectedRoute";

const UserRoutes = () => {
  return (

    <ProtectedRoute allowedRoles={["usuario"]} >
      <UserLayout>
        <Routes>
           
      
        </Routes>
      </UserLayout>
    </ProtectedRoute>
  );
};

export default UserRoutes;
