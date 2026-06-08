
import "../styles/mainLayout.css";
import NavigationMenu from "./NavigationMenu.jsx";
function MainLayout({ children }) {
    return (
        <main className="auth-page">
            <header className="auth-header">
                <div
                    className="logo-text"
                    aria-label="Чирик"
                >
                    ЧИРИК
                </div>
            </header>

            {children}

            <NavigationMenu />
        </main>
    );
}

export default MainLayout;