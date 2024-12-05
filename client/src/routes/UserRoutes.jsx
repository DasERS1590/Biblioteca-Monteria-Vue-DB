import React from "react";
import { Route, Routes } from "react-router-dom";

import UserLayout from "../components/Layout/UserLayout";
import ProtectedRoute from "./ProtectedRoute";
import Libro from "../components/User/Libro";
import Historial from "../components/User/Historial";

const UserRoutes = () => {
  return (

    <ProtectedRoute allowedRoles={["usuario"]} >
      <UserLayout>
        <Routes>
           
          <Route path="/libro" element={ <Libro/> } />
          <Route path="/historial" element={ <Historial/> } />
        </Routes>
      </UserLayout>
    </ProtectedRoute>
  );
};

export default UserRoutes;
