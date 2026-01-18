/**
 * Helper function to parse JSON safely
 */
const parseJSON = async (response) => {
  try {
    return await response.json();
  } catch {
    return null;
  }
};

/**
 * Login user
 */
export const loginUser = async (username, password) => {
  try {
    const response = await fetch("/login", {
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
 */
export const registerUser = async (username, password) => {
  try {
    const response = await fetch("/register", {
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
 * Fetch all links
 */
export const fetchLinks = async (token) => {
  try {
    const response = await fetch("/links", {
      method: "GET",
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    const data = await parseJSON(response);

    if (!response.ok) {
      console.error("Fetch links failed:", data);
      return [];
    }

    return Array.isArray(data) ? data : [];
  } catch (err) {
    console.error("Fetch links error:", err);
    return [];
  }
};

/**
 * Shorten URL
 */
export const shortenUrl = async (url, token) => {
  try {
    const response = await fetch("/shorten", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ url }),
    });

    const data = await parseJSON(response);

    if (!response.ok) {
      return { error: data?.message || "Failed to shorten URL" };
    }

    return data;
  } catch (err) {
    console.error("Shorten URL error:", err);
    return { error: "Network error" };
  }
};

/**
 * Delete URL
 */
export const deleteUrl = async (code, token) => {
  try {
    const response = await fetch("/delete", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ code }),
    });

    const data = await parseJSON(response);

    if (!response.ok) {
      return { error: data?.message || "Failed to delete URL" };
    }

    return data;
  } catch (err) {
    console.error("Delete URL error:", err);
    return { error: "Network error" };
  }
};
