import React from "react";
import { Route, Routes } from "react-router-dom";
import AdminDashboard from "../components/Admin/AdminDashboard";
import AdminItems from "../components/Admin/AdminItems";
import AdminSettings from "../components/Admin/AdminSettings";
import AdminLayout from "../components/Layout/AdminLayout";

const AdminRoutes = () => {
  return (
    <AdminLayout>
      <Routes>
        <Route path="/dashboard" element={<AdminDashboard />} />
        <Route path="/items" element={<AdminItems />} />
        <Route path="/settings" element={<AdminSettings />} />
      </Routes>
    </AdminLayout>
  );
};

export default AdminRoutes;
