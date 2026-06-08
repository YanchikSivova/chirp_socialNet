import {Link, useLocation} from 'react-router-dom';

import homeIcon from '../assets/home.svg';
import homeFilledIcon from '../assets/home-filled.svg';
import chatsIcon from "../assets/chats.svg";
import chatsFilledIcon from "../assets/chats-filled.svg";
import addIcon from "../assets/add.svg";
import bellIcon from "../assets/bell.svg";
import bellFilledIcon from "../assets/bell-filled.svg";
import userIcon from "../assets/user.svg"
import userFilledIcon from "../assets/user-filled.svg";


function NavigationMenu() {
  const location = useLocation();

  const currentPath = location.pathname;

  const menuItems = [
    { path: '/', icon: homeIcon, activeIcon: homeFilledIcon ,label: 'Главная' },
    { path: '/chats', icon: chatsIcon, activeIcon:chatsFilledIcon, label: 'Чаты' },
    { path: '/add', icon: addIcon, activeIcon: addIcon, label: 'Добавить' },
    { path: '/notifications', icon: bellIcon, activeIcon: bellFilledIcon, label: 'Уведомления' },
    { path: '/@me', icon: userIcon, activeIcon:userFilledIcon, label: 'Профиль' },
  ];

  const isActive = (path) => {
    if (path === '/') {
      return currentPath === '/';
    }

    if (path === '/@me') {
      return currentPath.startsWith('/@');
    }

    return currentPath === path;
  };

  return (
      <nav className="nav-menu">
        <ul className="nav-menu-list">
          {menuItems.map((item) => {
            const active = isActive(item.path);

            return (
                <li
                    key={item.path}
                    className={`nav-menu-item ${
                        active ? 'active-icon' : 'inactive-icon'
                    }`}
                >
                  <Link
                      to={item.path}
                      className={`nav-menu-link${active ? ' active' : ''}`}
                      aria-label={item.label}
                  >
                    <img
                        src={active ? item.activeIcon : item.icon}
                        alt={item.label}
                        className="nav-menu-icon"
                    />
                  </Link>
                </li>
            );
          })}
        </ul>
      </nav>
  );
}

export default NavigationMenu;
