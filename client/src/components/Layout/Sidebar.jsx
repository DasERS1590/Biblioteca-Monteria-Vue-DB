import { NavLink, useNavigate } from "react-router-dom";
import "../../styles/common/Sidebar.css";

const Sidebar = ({ links = [], logo = "Biblioteca", expanded = true, onToggle }) => {
  const navigate = useNavigate();
  const user = JSON.parse(localStorage.getItem("user"));

  const handleLogout = () => {
    localStorage.removeItem("user");
    navigate("/login");
  };

  return (
    <aside className={`sidebar ${expanded ? 'expanded' : 'collapsed'}`}>
      <div className="sidebar-header">
        <div className="sidebar-logo">{expanded ? logo : logo.charAt(0)}</div>
        <button 
          className="sidebar-toggle-btn"
          onClick={onToggle}
          title={expanded ? "Contraer menú" : "Expandir menú"}
        >
          {expanded ? '◀' : '▶'}
        </button>
      </div>
      
      <nav className="sidebar-links">
        {links.map((link) => (
          <NavLink
            key={link.to}
            to={link.to}
            className={({ isActive }) =>
              isActive ? "sidebar-link active" : "sidebar-link"
            }
            end
            title={!expanded ? link.label : ""}
          >
            {expanded ? link.label : link.label.charAt(0)}
          </NavLink>
        ))}
      </nav>
      
      <div className="sidebar-user-logout">
        <div className="sidebar-user" title={!expanded ? `${user?.rol}: ${user?.nombre}` : ""}>
          {expanded ? `${user?.rol}: ${user?.nombre}` : `${user?.rol?.charAt(0)}: ${user?.nombre?.charAt(0)}`}
        </div>
        <button 
          className="sidebar-logout" 
          onClick={handleLogout}
          title={!expanded ? "Cerrar sesión" : ""}
        >
          {expanded ? "Cerrar sesión" : "🚪"}
        </button>
      </div>
    </aside>
  );
};

export default Sidebar; 