import { initializeApp, getApp, getApps, type FirebaseApp } from "firebase/app";
import { getAuth, GoogleAuthProvider } from "firebase/auth";

const firebaseConfig = {
    apiKey: "AIzaSyBu1opPNUjh7kqCjGJXvzTYc4g1aWQE0T0",
    authDomain: "fif-tax.firebaseapp.com",
    projectId: "fif-tax",
    storageBucket: "fif-tax.firebasestorage.app",
    messagingSenderId: "648250511998",
    appId: "1:648250511998:web:b5a6d1db22b87512c0be2f",
    measurementId: "G-P3VF39S59K",
};

let app: FirebaseApp;
if (!getApps().length) {
    app = initializeApp(firebaseConfig);
} else {
    app = getApp();
}

export const auth = getAuth(app);
export const googleProvider = new GoogleAuthProvider();
