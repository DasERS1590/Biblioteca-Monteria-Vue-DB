import { Route, Routes } from "react-router-dom";
import Dashboard from "../components/Admin/Dashboard";
import Books from "../components/Admin/Books";
import Users from "../components/Admin/Users";
import AdminLayout from "../components/Layout/AdminLayout";
import ProtectedRoute from "./ProtectedRoute";
import Loans from "../components/Admin/Loans";
import Fines from "../components/Admin/Fines";
import Reservation from "../components/Admin/Reservation";
import Register from "../pages/Register";
import CreateBook from "../components/Admin/CreateBook";
import CreateEditorial from "../components/Admin/CreateEditorial";
import CreateAuthor  from "../components/Admin/Author";


const AdminRoutes = () => {
  return (
    <ProtectedRoute allowedRoles={["administrador"]}>
      <AdminLayout>
        <Routes>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/books" element={<Books/>} />
          <Route path="/users" element={<Users/>} />
          <Route path="/loans" element={<Loans/>} />
          <Route path="/fines" element={<Fines/>} />
          <Route path="/reservations" element={<Reservation/>} />
          <Route path="/registrarlibro" element={ <CreateBook/> }/>
          <Route path="/registaredi" element = {<CreateEditorial/>} />
          <Route path="/registrar-autor" element ={<CreateAuthor/>} />
        </Routes>
      </AdminLayout>
      </ProtectedRoute >

  );
};

export default AdminRoutes;
