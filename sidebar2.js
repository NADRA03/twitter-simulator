export function render() {
    return `
        <div class="sidebar-section">
            <h2>Who to Follow</h2>
            <ul>
                <li><a href="/user/1" onclick="loadPage('user'); return false;">User 1</a></li>
                <li><a href="/user/2" onclick="loadPage('user'); return false;">User 2</a></li>
                <li><a href="/user/3" onclick="loadPage('user'); return false;">User 3</a></li>
                <li><a href="/user/4" onclick="loadPage('user'); return false;">User 4</a></li>
                <!-- Add more users as needed -->
            </ul>
        </div>
        <div class="sidebar-section">
            <h2>Trends for You</h2>
            <ul>
                <li><a href="/trend/1" onclick="loadPage('trend'); return false;">Trend 1</a></li>
                <li><a href="/trend/2" onclick="loadPage('trend'); return false;">Trend 2</a></li>
                <li><a href="/trend/3" onclick="loadPage('trend'); return false;">Trend 3</a></li>
                <li><a href="/trend/4" onclick="loadPage('trend'); return false;">Trend 4</a></li>
                <!-- Add more trends as needed -->
            </ul>
        </div>
    `;
}