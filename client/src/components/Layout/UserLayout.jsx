import { useState } from "react";
import Sidebar from "./Sidebar";
import Breadcrumbs from "../common/Breadcrumbs";

const userLinks = [
  { to: "/user/dashboard", label: "Dashboard" },
  { to: "/user/libro", label: "Libros" },
  { to: "/user/prestamo", label: "Préstamos" },
  { to: "/user/reserva", label: "Reservas" },
  { to: "/user/multa", label: "Multas" },
  { to: "/user/historial", label: "Historial" },
];

const UserLayout = ({ children }) => {
  const [sidebarExpanded, setSidebarExpanded] = useState(true);

  const toggleSidebar = () => {
    setSidebarExpanded(!sidebarExpanded);
  };

  return (
    <div className={`layout ${sidebarExpanded ? 'sidebar-expanded' : 'sidebar-collapsed'}`}>
      <Sidebar 
        links={userLinks} 
        logo="Biblioteca" 
        expanded={sidebarExpanded}
        onToggle={toggleSidebar}
      />
      <div className="content">
        <Breadcrumbs />
        {children}
      </div>
    </div>
  );
};

export default UserLayout;
