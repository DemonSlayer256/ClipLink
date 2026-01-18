import { createContext, useState } from 'react';

export const AuthContext = createContext();

export const AuthProvider = ({ children }) => {
    const [isAuth, setIsAuth] = useState(!!localStorage.getItem('token'));

    return (
        <AuthContext.Provider value={{ isAuth, setAuth: setIsAuth }}>
            {children}
        </AuthContext.Provider>
    );
};
