import React from "react";
import { Route, Routes } from "react-router-dom";
import UserDashboard from "../components/User/UserDashboard";
import UserOrders from "../components/User/UserOrders";
import UserProfile from "../components/User/UserProfile";
import UserLayout from "../components/Layout/UserLayout";
import ProtectedRoute from "./ProtectedRoute";

const UserRoutes = () => {
  return (

    <ProtectedRoute allowedRoles={["usuario"]} >
      <UserLayout>
        <Routes>
          <Route path="/dashboard" element={<UserDashboard />} />
          <Route path="/orders" element={<UserOrders />} />
          <Route path="/profile" element={<UserProfile />} />
        </Routes>
      </UserLayout>
    </ProtectedRoute>
  );
};

export default UserRoutes;
