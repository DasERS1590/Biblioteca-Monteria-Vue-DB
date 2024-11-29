import React from "react";
import "../../styles/Admin.css";

const AdminItems = () => {
  return (
    <div className="admin-items">
      <h1>Items Management</h1>
      <button className="add-item-btn">+ Add Item</button>
      <table className="items-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Title</th>
            <th>Description</th>
            <th>Actions</th>
          </tr>
        </thead>
        <tbody>
          <tr>
            <td>1</td>
            <td>Item #1</td>
            <td>Item description</td>
            <td>
              <button className="action-btn">...</button>
            </td>
          </tr>
          <tr>
            <td>2</td>
            <td>Item #2</td>
            <td>Item description</td>
            <td>
              <button className="action-btn">...</button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );
};

export default AdminItems;
