const API_URL = import.meta.env.VITE_API_URL;

/**
 * Helper function to parse JSON safely
 */
const parseJSON = async (response) => {
  try {
    const data = await response.json();
    return data;
  } catch (err) {
    return null;
  }
};

/**
 * Login user
 * Returns token string if success, null otherwise
 */
export const loginUser = async (username, password) => {
  try {
    const response = await fetch(`${API_URL}/login`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ user: username, pass: password }),
    });

    const data = await parseJSON(response);

    if (!response.ok || !data?.token) {
      console.error("Login failed:", data);
      return null;
    }

    return data.token;
  } catch (err) {
    console.error("Login error:", err);
    return null;
  }
};

/**
 * Register user
 * Returns server response (success message or error)
 */
export const registerUser = async (username, password) => {
  try {
    const response = await fetch(`${API_URL}/register`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ user: username, pass: password }),
    });

    const data = await parseJSON(response);
    if (!response.ok) {
      console.error("Register failed:", data);
    }
    return data;
  } catch (err) {
    console.error("Register error:", err);
    return { error: "Network error" };
  }
};

/**
 * Fetch all links for the user
 * Always returns an array of URLMapping objects
 */
export const fetchLinks = async (token) => {
  try {
    const response = await fetch(`${API_URL}/links`, {
      method: "GET",
      headers: { Authorization: `Bearer ${token}` },
    });

    const data = await parseJSON(response);
    if (!response.ok) {
      console.error("Fetch links failed:", data);
      return [];
    }

    // If API returns an array directly
    if (Array.isArray(data)) return data;

    // If API returns object { links: [...] } or similar
    if (Array.isArray(data.links)) return data.links;

    // Fallback: single object or invalid data
    return [];
  } catch (err) {
    console.error("Fetch links error:", err);
    return [];
  }
};

/**
 * Shorten a URL
 * Returns a URLMapping object on success or { error: ... } on failure
 */
export const shortenUrl = async (url, token) => {
  try {
    const response = await fetch(`${API_URL}/shorten`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ url }),
    });

    const data = await parseJSON(response);

    if (!response.ok) {
      console.error("Shorten URL failed:", data);
      return { error: data?.message || "Failed to shorten URL" };
    }

    return data;
  } catch (err) {
    console.error("Shorten URL error:", err);
    return { error: "Network error" };
  }
};


export const deleteUrl = async (code, token) => {
  try {
    const response = await fetch(`${API_URL}/delete`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ code }),
    });

    const data = await parseJSON(response);

    if (!response.ok) {
      console.error("Delete URL failed:", response);
      return { error: data?.message || "Failed to delete URL" };
    }

    return data;
  } catch (err) {
    console.error("Delete URL error:", err);
    return { error: "Network error" };
  }
};
