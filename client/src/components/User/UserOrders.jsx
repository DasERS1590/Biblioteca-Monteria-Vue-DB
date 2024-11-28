import React, { useState } from "react";
import "../../styles/User.css";

const UserOrders = () => {
  const [orders, setOrders] = useState([
    {
      id: 1,
      date: "2024-11-25",
      total: "$50.00",
      status: "Delivered",
    },
    {
      id: 2,
      date: "2024-11-20",
      total: "$30.00",
      status: "Processing",
    },
    {
      id: 3,
      date: "2024-11-15",
      total: "$70.00",
      status: "Cancelled",
    },
  ]);

  return (
    <div className="content">
      <h2>User Orders</h2>
      <div className="card">
        <h3>Order History</h3>
        <table>
          <thead>
            <tr>
              <th>Order ID</th>
              <th>Date</th>
              <th>Total</th>
              <th>Status</th>
            </tr>
          </thead>
          <tbody>
            {orders.map((order) => (
              <tr key={order.id}>
                <td>{order.id}</td>
                <td>{order.date}</td>
                <td>{order.total}</td>
                <td>{order.status}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default UserOrders;
