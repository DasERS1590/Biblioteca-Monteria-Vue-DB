import { useState } from "react";
import Sidebar from "./Sidebar";
import Breadcrumbs from "../common/Breadcrumbs";
import "../../styles/common/Admin.css";

const adminLinks = [
  { to: "/admin/dashboard", label: "Dashboard" },
  { to: "/admin/books", label: "Libros" },
  { to: "/admin/users", label: "Usuarios" },
  { to: "/admin/loans", label: "Préstamos" },
  { to: "/admin/fines", label: "Multas" },
  { to: "/admin/reservations", label: "Reservas" },
  { to: "/admin/registrarlibro", label: "Registrar libro" },
  { to: "/admin/registaredi", label: "Registrar Editorial" },
  { to: "/admin/registrar-autor", label: "Registrar Autor" },
];

const AdminLayout = ({ children }) => {
  const [sidebarExpanded, setSidebarExpanded] = useState(true);

  const toggleSidebar = () => {
    setSidebarExpanded(!sidebarExpanded);
  };

  return (
    <div className={`layout ${sidebarExpanded ? 'sidebar-expanded' : 'sidebar-collapsed'}`}>
      <Sidebar 
        links={adminLinks} 
        logo="Admin" 
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

export default AdminLayout;
