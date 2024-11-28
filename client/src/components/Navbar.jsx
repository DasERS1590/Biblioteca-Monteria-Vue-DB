import React from "react";
import { Link } from "react-router-dom";
import "../styles/Navar.css"

const Navbar = ({ isAdmin }) => {
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
        <li><Link to="/login">Logout</Link></li>
      </ul>
    </nav>
  );
};

export default Navbar;
