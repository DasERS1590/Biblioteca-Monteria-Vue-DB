import React, { useState } from "react";
import "../../styles/User.css";

const UserProfile = () => {
  const [user, setUser] = useState({
    fullName: "John Doe",
    email: "user@example.com",
    phone: "123-456-7890",
  });

  const [editMode, setEditMode] = useState(false);
  const [formData, setFormData] = useState(user);

  const handleInputChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const saveChanges = () => {
    setUser(formData);
    setEditMode(false);
  };

  return (
    <div className="content">
      <h2>User Profile</h2>
      {!editMode ? (
        <div className="card">
          <h3>Profile Information</h3>
          <p><strong>Full Name:</strong> {user.fullName}</p>
          <p><strong>Email:</strong> {user.email}</p>
          <p><strong>Phone:</strong> {user.phone}</p>
          <button onClick={() => setEditMode(true)}>Edit Profile</button>
        </div>
      ) : (
        <div className="card">
          <h3>Edit Profile</h3>
          <form>
            <div>
              <label>Full Name</label>
              <input
                type="text"
                name="fullName"
                value={formData.fullName}
                onChange={handleInputChange}
              />
            </div>
            <div>
              <label>Email</label>
              <input
                type="email"
                name="email"
                value={formData.email}
                onChange={handleInputChange}
              />
            </div>
            <div>
              <label>Phone</label>
              <input
                type="text"
                name="phone"
                value={formData.phone}
                onChange={handleInputChange}
              />
            </div>
            <button type="button" onClick={saveChanges}>Save</button>
            <button type="button" onClick={() => setEditMode(false)}>Cancel</button>
          </form>
        </div>
      )}
    </div>
  );
};

export default UserProfile;
