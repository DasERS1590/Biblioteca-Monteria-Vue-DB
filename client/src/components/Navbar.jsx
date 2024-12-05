import React from "react";
import "../styles/Navar.css"
import { Link, useNavigate } from "react-router-dom";

const Navbar = ({ isAdmin }) => {

  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem("user");
    console.log("Usuario deslogueado");
    navigate("/login");

  };

  return (
    <nav className="navbar">
      <h2>{isAdmin ? "Panel Adminstrador" : "Panel de usuario"}</h2>
      <ul>
        {isAdmin ? (
          <>
    
            <li><Link to="/admin/books">Libros</Link></li>
            <li><Link to="/admin/users">Usuarios</Link></li>
            <li><Link to="/admin/loans">Préstamos</Link></li>
            <li><Link to="/admin/fines">Multas</Link></li>
            <li><Link to="/admin/reservations">Reservas</Link></li>
            <li><Link to="/admin/register">Registar Usuario </Link></li>
            

          </>
        ) : (
          <>
            <li><Link to="/user/dashboard">Dashboard</Link></li>
          
          </>
        )}
        <button onClick={handleLogout} style={{ cursor: "pointer", background: "none", border: "none", color: "blue", textDecoration: "underline" }}>
          Logout
        </button>
      </ul>
    </nav>
  );
};

export default Navbar;
