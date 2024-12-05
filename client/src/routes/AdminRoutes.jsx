import React from "react";
import { Route, Routes } from "react-router-dom";
import Books from "../components/Admin/Books";
import Users from "../components/Admin/Users";
import AdminLayout from "../components/Layout/AdminLayout";
import ProtectedRoute from "./ProtectedRoute";
import Loans from "../components/Admin/Loans";
import Fines from "../components/Admin/Fines";
import Reservation from "../components/Admin/Reservation";
import Register from "../components/Admin/Register";


const AdminRoutes = () => {
  return (
    <ProtectedRoute allowedRoles={["administrador"]}>
      <AdminLayout>
        <Routes>
          <Route path="/books" element={<Books/>} />
          <Route path="/users" element={<Users/>} />
          <Route path="/loans" element={<Loans/>} />
          <Route path="/fines" element={<Fines/>} />
          <Route path="/reservations" element={<Reservation/>} />
          <Route path="/register" element={ <Register/> }/>
        </Routes>
      </AdminLayout>
      </ProtectedRoute >

  );
};

export default AdminRoutes;
