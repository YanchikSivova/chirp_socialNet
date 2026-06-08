import { useState } from 'react';

import {
    searchByName,
    searchByUsername,
} from '../api/client';

import NavigationMenu from "../components/NavigationMenu.jsx";
import { useNavigate } from 'react-router-dom';
import "../styles/search.css";

import searchIcon from "../assets/search.svg";

function SearchPage({ me }) {
    const [query, setQuery] = useState('');
    const [mode, setMode] = useState('username');
    const [profiles, setProfiles] = useState([]);
    const [loading, setLoading] = useState(false);
    const navigate = useNavigate();
    async function handleSearch() {
        if (!query.trim()) return;

        setLoading(true);

        try {
            let data;

            if (mode === 'name') {
                data = await searchByName(query);
            } else {
                data = await searchByUsername(query);
            }

            setProfiles(data?.profiles || []);
        } finally {
            setLoading(false);
        }
    }

    return (
        <div className="search-page">
            <div className="search-input-wrapper">
                <input
                    value={query}
                    onChange={(e) =>
                        setQuery(e.target.value)
                    }
                    placeholder="Поиск..."
                    className="search-input"
                    onKeyDown={(e) => {
                        if (e.key === 'Enter') {
                            handleSearch();
                        }
                    }}
                />

                <button
                    className="search-submit-btn"
                    onClick={handleSearch}
                >
                    <img
                        src={searchIcon}
                        alt="Поиск"
                    />
                </button>
            </div>

            <div className="search-tabs">
                <button
                    className={`search-tab ${
                        mode === 'username'
                            ? 'active'
                            : ''
                    }`}
                    onClick={() =>
                        setMode('username')
                    }
                >
                    По юзернейму
                </button>

                <button
                    className={`search-tab ${
                        mode === 'name'
                            ? 'active'
                            : ''
                    }`}
                    onClick={() =>
                        setMode('name')
                    }
                >
                    По имени
                </button>
            </div>

            {loading && (
                <p className="search-loading">
                    Поиск...
                </p>
            )}

            <div className="search-results">
                {profiles.map((p) => (
                    <div
                        key={p.profile_id}
                        className="search-profile-item"
                        onClick={() => {
                            if (me?.profile_id === p.profile_id) {
                                navigate('/@me');
                            } else {
                                navigate(`/@${p.profile_id}`);
                            }
                        }}
                    >
                        {p.avatar ? (
                            <img
                                src={p.avatar}
                                alt={p.name}
                            />
                        ) : (
                            <div className="avatar-fallback">
                                {p.name?.[0]}
                            </div>
                        )}

                        <div className="search-profile-info">
                            <strong>{p.name}</strong>

                            <span>@{p.username}</span>
                        </div>
                    </div>
                ))}
            </div>

            <NavigationMenu />
        </div>
    );
}

export default SearchPage;