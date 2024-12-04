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
      <h2>{isAdmin ? "Admin Panel" : "User Panel"}</h2>
      <ul>
        {isAdmin ? (
          <>
            <li><Link to="/admin/dashboard">Dashboard</Link></li>
            <li><Link to="/admin/items">Items</Link></li>
            <li><Link to="/admin/settings">Settings</Link></li>
          </>
        ) : (
          <>
            <li><Link to="/user/dashboard">Dashboard</Link></li>
            <li><Link to="/user/orders">Orders</Link></li>
            <li><Link to="/user/profile">Profile</Link></li>
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
