const API_URL = 'http://localhost:8080';

export const albumService = {
  // Get all albums
  getAll: async () => {
    const response = await fetch(`${API_URL}/albums`);
    if (!response.ok) {
      throw new Error('Failed to fetch albums');
    }
    return response.json();
  },

  // Get album by ID
  getById: async (id) => {
    const response = await fetch(`${API_URL}/albums/${id}`);
    if (!response.ok) {
      throw new Error('Album not found');
    }
    return response.json();
  },

  // Create new album
  create: async (album) => {
    const response = await fetch(`${API_URL}/albums`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(album),
    });
    if (!response.ok) {
      throw new Error('Failed to create album');
    }
    return response.json();
  },

  // Update album
  update: async (id, album) => {
    const response = await fetch(`${API_URL}/albums/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(album),
    });
    if (!response.ok) {
      throw new Error('Failed to update album');
    }
    return response.json();
  },

  // Delete album
  delete: async (id) => {
    const response = await fetch(`${API_URL}/albums/${id}`, {
      method: 'DELETE',
    });
    if (!response.ok) {
      throw new Error('Failed to delete album');
    }
    return response.json();
  },
};
