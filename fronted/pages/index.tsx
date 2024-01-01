import React, { useState, useEffect } from 'react';
import './_app'; // Import your style file

const UserOperationsPage = () => {
  const [users, setUsers] = useState([]);
  const [selectedUser, setSelectedUser] = useState({ id: '', username: '', email: '' });
  const [newUser, setNewUser] = useState({ username: '', email: '' });

  const fetchUsers = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/users');
      const data = await response.json();
      setUsers(data);
    } catch (error) {
      console.error('Error fetching users:', error);
    } 
  };

  useEffect(() => {
    // Fetching users
    fetchUsers();
  }, []);

  const handleUserSelect = (user) => {
    // Selecting a user
    setSelectedUser(user);
  };

  const handleUpdateUser = async () => {
    // Updating a user
    try {
      const response = await fetch(`http://localhost:8080/api/users/${selectedUser.id}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(selectedUser),
      });

      if (response.ok) {
        // If update is successful, fetch users again
        fetchUsers();
      } else {
        console.error('Error updating user:', response.statusText);
      }
    } catch (error) {
      console.error('Error updating user:', error);
    }
  };

  const handleDeleteUser = async () => {
    // Deleting a user
    try {
      const response = await fetch(`http://localhost:8080/api/users/${selectedUser.id}`, {
        method: 'DELETE',
      });

      if (response.ok) {
        // If deletion is successful, fetch users again
        await response.json(); // Added this line
        fetchUsers();
        setSelectedUser({ id: '', username: '', email: '' });
      } else {
        console.error('Error deleting user:', response.statusText);
      }
    } catch (error) {
      console.error('Error deleting user:', error);
    }
  };

  const handleAddUser = async () => {
    // Adding a new user
    try {
      const response = await fetch('http://localhost:8080/api/users', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(newUser),
      });

      if (response.ok) {
        // If addition is successful, fetch users again
        fetchUsers();
        setNewUser({ username: '', email: '' });
      } else {
        console.error('Error adding user:', response.statusText);
      }
    } catch (error) {
      console.error('Error adding user:', error);
    }
  };

  return (
    <div className="user-operations-container">
      <h1>User Operations</h1>

      {/* Form for adding a new user */}
      <h2>Add New User</h2>
      <label>Username:</label>
      <input
        type="text"
        value={newUser.username}
        onChange={(e) => setNewUser({ ...newUser, username: e.target.value })}
      />
      <label>Email:</label>
      <input
        type="text"
        value={newUser.email}
        onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
      />
      <button onClick={handleAddUser}>Add</button>

      {/* List of users */}
      <h2>User List</h2>
      {users ? (
        <ul>
          {users.map((user) => (
            <li key={user.id} onClick={() => handleUserSelect(user)}>
              {user.username} - {user.email}
            </li>
          ))}
        </ul>
      ) : (
        <p>No users found.</p>
      )}

      {/* Details of the selected user */}
      {selectedUser.id && (
        <div>
          <h2>Selected User Details</h2>
          <label>ID:</label>
          <input type="text" value={selectedUser.id} readOnly />
          <label>Username:</label>
          <input
            type="text"
            value={selectedUser.username}
            onChange={(e) => setSelectedUser({ ...selectedUser, username: e.target.value })}
          />
          <label>Email:</label>
          <input
            type="text"
            value={selectedUser.email}
            onChange={(e) => setSelectedUser({ ...selectedUser, email: e.target.value })}
          />
          <button onClick={handleUpdateUser}>Update</button>
          <button onClick={handleDeleteUser}>Delete</button>
        </div>
      )}
    </div>
  );
};

export default UserOperationsPage;
