import { Route, Routes } from "react-router-dom";

import UserLayout from "../components/Layout/UserLayout";
import ProtectedRoute from "./ProtectedRoute";
import Dashboard from "../components/User/Dashboard";
import Libro from "../components/User/Libro";
import Historial from "../components/User/Historial";
import Multa from "../components/User/Multas";
import Prestamo from "../components/User/Prestamos";
import Reserva from "../components/User/Reservas";

const UserRoutes = () => {
  return (

    <ProtectedRoute allowedRoles={["usuario"]} >
      <UserLayout>
        <Routes>
          <Route path="/dashboard" element={<Dashboard />} />
          <Route path="/libro" element={ <Libro/> } />
          <Route path="/historial" element={ <Historial/> } />
          <Route path="/multa" element = {<Multa/>} />
          <Route path="/prestamo" element = { <Prestamo/>} />
          <Route path="/reserva" element = { <Reserva/> }/>

        </Routes>
      </UserLayout>
    </ProtectedRoute>
  );
};

export default UserRoutes;
