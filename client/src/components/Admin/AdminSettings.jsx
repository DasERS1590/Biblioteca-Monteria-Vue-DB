import React from "react";
import "../../styles/Admin.css";

const AdminSettings = () => {
  return (
    <div className="admin-settings">
      <h1>User Settings</h1>
      <div className="settings-tabs">
        <button className="tab active">My profile</button>
        <button className="tab">Password</button>
        <button className="tab">Appearance</button>
      </div>
      <form className="settings-form">
        <label>
          Full Name
          <input type="text" defaultValue="Admin" />
        </label>
        <label>
          Email
          <input type="email" defaultValue="admin@example.com" />
        </label>
        <button className="save-btn">Save</button>
        <button className="cancel-btn">Cancel</button>
      </form>
    </div>
  );
};

export default AdminSettings;
