#!/usr/bin/env python3
"""T3 Trip Planner MCP Server — serves a Mapbox UI with TSP trip results per day."""

import asyncio
import json
import os

import httpx
import mcp.types as types
from mcp.server import NotificationOptions, Server
from mcp.server.models import InitializationOptions
from mcp.server.stdio import stdio_server

BACKEND_URL = os.getenv("T3_BACKEND_URL", "http://localhost:8080")
MAPBOX_TOKEN = os.getenv("MAPBOX_TOKEN", "")

DAY_COLORS = [
    "#e74c3c", "#3498db", "#2ecc71", "#f39c12",
    "#9b59b6", "#1abc9c", "#e91e63", "#ff5722",
]

server = Server("t3-trip-planner")


def _build_html(data: dict) -> str:
    location = data.get("Location", "Unknown")
    days_count = data.get("Days", 0)
    tsp: list[list[dict]] = data.get("TSP", [])

    all_sites = [s for day in tsp for s in day]
    if all_sites:
        center_lat = sum(s["lat"] for s in all_sites) / len(all_sites)
        center_lng = sum(s["lng"] for s in all_sites) / len(all_sites)
    else:
        center_lat, center_lng = 46.6034, 1.8883  # France centre

    sidebar_days = ""
    for i, day in enumerate(tsp):
        color = DAY_COLORS[i % len(DAY_COLORS)]
        sites_html = ""
        for j, site in enumerate(day):
            name = site.get("name", "—")
            stype = site.get("type", "")
            city = site.get("city", "")
            meta = " · ".join(filter(None, [stype, city]))
            sites_html += f"""
            <div class="site-item" onclick="focusSite({site['lng']}, {site['lat']}, {i}, {j})">
              <span class="site-num" style="background:{color}">{j + 1}</span>
              <div>
                <div class="site-name">{name}</div>
                <div class="site-meta">{meta}</div>
              </div>
            </div>"""

        sidebar_days += f"""
        <div class="day-block">
          <div class="day-header" style="border-left:4px solid {color}" onclick="focusDay({i})">
            <span>Day {i + 1}</span>
            <span class="day-count">{len(day)} sites</span>
          </div>
          {sites_html}
        </div>"""

    tsp_json = json.dumps(tsp)
    colors_json = json.dumps(DAY_COLORS)

    return f"""<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{location} — {days_count}-Day Trip</title>
  <script src="https://api.mapbox.com/mapbox-gl-js/v3.3.0/mapbox-gl.js"></script>
  <link href="https://api.mapbox.com/mapbox-gl-js/v3.3.0/mapbox-gl.css" rel="stylesheet">
  <style>
    *, *::before, *::after {{ box-sizing: border-box; margin: 0; padding: 0; }}
    body {{ font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
           display: flex; height: 100vh; overflow: hidden; background: #f5f5f5; }}

    /* ── Sidebar ── */
    #sidebar {{
      width: 320px; min-width: 320px; display: flex; flex-direction: column;
      background: #fff; box-shadow: 2px 0 16px rgba(0,0,0,.12); z-index: 1; overflow: hidden;
    }}
    .sidebar-header {{
      padding: 18px 16px 14px; background: #0f172a; color: #fff; flex-shrink: 0;
    }}
    .sidebar-header h1 {{ font-size: 17px; font-weight: 700; letter-spacing: -.3px; }}
    .sidebar-header p  {{ font-size: 12px; color: #94a3b8; margin-top: 3px; }}
    #sidebar-body {{ overflow-y: auto; flex: 1; }}

    .day-block {{ border-bottom: 1px solid #f1f5f9; }}
    .day-header {{
      display: flex; justify-content: space-between; align-items: center;
      padding: 11px 14px; cursor: pointer; font-size: 13px; font-weight: 600;
      color: #1e293b; background: #f8fafc; transition: background .15s;
    }}
    .day-header:hover {{ background: #f1f5f9; }}
    .day-count {{ font-size: 11px; font-weight: 400; color: #64748b; }}

    .site-item {{
      display: flex; align-items: flex-start; gap: 10px;
      padding: 9px 14px 9px 18px; cursor: pointer; transition: background .15s;
    }}
    .site-item:hover {{ background: #f8fafc; }}
    .site-item.active {{ background: #eff6ff; }}
    .site-num {{
      flex-shrink: 0; width: 22px; height: 22px; border-radius: 50%;
      color: #fff; font-size: 11px; font-weight: 700;
      display: flex; align-items: center; justify-content: center; margin-top: 1px;
    }}
    .site-name {{ font-size: 13px; font-weight: 500; color: #1e293b; }}
    .site-meta {{ font-size: 11px; color: #94a3b8; margin-top: 2px; text-transform: capitalize; }}

    /* ── Map ── */
    #map {{ flex: 1; }}

    /* ── Popups ── */
    .mapboxgl-popup-content {{
      border-radius: 10px; padding: 0; overflow: hidden;
      box-shadow: 0 8px 24px rgba(0,0,0,.15); min-width: 200px;
    }}
    .popup-header {{
      padding: 10px 14px 8px; color: #fff; font-size: 12px; font-weight: 600;
    }}
    .popup-body {{ padding: 10px 14px 12px; }}
    .popup-title {{ font-size: 14px; font-weight: 700; color: #0f172a; margin-bottom: 3px; }}
    .popup-type  {{ font-size: 12px; color: #64748b; text-transform: capitalize; }}
    .popup-city  {{ font-size: 11px; color: #94a3b8; }}
    .popup-desc  {{ font-size: 12px; color: #475569; margin-top: 7px; line-height: 1.45; }}
    .popup-link  {{
      display: inline-block; margin-top: 8px; font-size: 12px;
      color: #3b82f6; text-decoration: none; font-weight: 500;
    }}
    .popup-link:hover {{ text-decoration: underline; }}
  </style>
</head>
<body>
  <div id="sidebar">
    <div class="sidebar-header">
      <h1>{location}</h1>
      <p>{days_count}-day itinerary &middot; {len(all_sites)} sites</p>
    </div>
    <div id="sidebar-body">{sidebar_days}</div>
  </div>
  <div id="map"></div>

  <script>
    mapboxgl.accessToken = '{MAPBOX_TOKEN}';

    const tsp    = {tsp_json};
    const colors = {colors_json};
    let activePopup = null;

    const map = new mapboxgl.Map({{
      container: 'map',
      style: 'mapbox://styles/mapbox/streets-v12',
      center: [{center_lng}, {center_lat}],
      zoom: 11
    }});

    map.addControl(new mapboxgl.NavigationControl(), 'top-right');
    map.addControl(new mapboxgl.ScaleControl({{ maxWidth: 120, unit: 'metric' }}));

    map.on('load', () => {{
      tsp.forEach((day, di) => {{
        const color = colors[di % colors.length];

        // Route line
        if (day.length > 1) {{
          const coords = day.map(s => [s.lng, s.lat]);
          map.addSource(`route-${{di}}`, {{
            type: 'geojson',
            data: {{ type: 'Feature', geometry: {{ type: 'LineString', coordinates: coords }} }}
          }});
          map.addLayer({{
            id: `route-${{di}}`,
            type: 'line',
            source: `route-${{di}}`,
            paint: {{ 'line-color': color, 'line-width': 3, 'line-opacity': 0.7 }}
          }});
        }}

        // Markers
        day.forEach((site, si) => {{
          const el = document.createElement('div');
          Object.assign(el.style, {{
            width: '26px', height: '26px', borderRadius: '50%',
            background: color, border: '2.5px solid #fff',
            boxShadow: '0 2px 8px rgba(0,0,0,.28)', cursor: 'pointer',
            display: 'flex', alignItems: 'center', justifyContent: 'center',
            fontSize: '11px', fontWeight: '700', color: '#fff', fontFamily: 'sans-serif'
          }});
          el.textContent = si + 1;

          const descHtml  = site.description ? `<div class="popup-desc">${{site.description}}</div>` : '';
          const linkHtml  = site.website     ? `<a class="popup-link" href="${{site.website}}" target="_blank">Visit website →</a>` : '';
          const cityMeta  = [site.type, site.city].filter(Boolean).join(' · ');

          const popup = new mapboxgl.Popup({{ offset: 16, maxWidth: '260px', closeButton: false }})
            .setHTML(`
              <div class="popup-header" style="background:${{color}}">Day ${{di + 1}} — stop ${{si + 1}}</div>
              <div class="popup-body">
                <div class="popup-title">${{site.name}}</div>
                <div class="popup-type">${{site.type}}</div>
                <div class="popup-city">${{site.city || ''}}</div>
                ${{descHtml}}
                ${{linkHtml}}
              </div>`);

          new mapboxgl.Marker(el)
            .setLngLat([site.lng, site.lat])
            .setPopup(popup)
            .addTo(map);
        }});
      }});

      fitAll();
    }});

    function fitAll() {{
      const coords = tsp.flat().map(s => [s.lng, s.lat]);
      if (!coords.length) return;
      const bounds = coords.reduce((b, c) => b.extend(c),
        new mapboxgl.LngLatBounds(coords[0], coords[0]));
      map.fitBounds(bounds, {{ padding: 60, maxZoom: 14 }});
    }}

    function focusDay(di) {{
      const day = tsp[di];
      if (!day?.length) return;
      const bounds = day.reduce((b, s) => b.extend([s.lng, s.lat]),
        new mapboxgl.LngLatBounds([day[0].lng, day[0].lat], [day[0].lng, day[0].lat]));
      map.fitBounds(bounds, {{ padding: 80, maxZoom: 14, duration: 900 }});

      tsp.forEach((_, i) => {{
        if (map.getLayer(`route-${{i}}`)) {{
          map.setPaintProperty(`route-${{i}}`, 'line-opacity', i === di ? 1 : 0.15);
          map.setPaintProperty(`route-${{i}}`, 'line-width',   i === di ? 5 : 1.5);
        }}
      }});
    }}

    function focusSite(lng, lat, di, si) {{
      focusDay(di);
      map.flyTo({{ center: [lng, lat], zoom: 15, duration: 800 }});
      document.querySelectorAll('.site-item').forEach(el => el.classList.remove('active'));
      const items = document.querySelectorAll('.site-item');
      let idx = 0;
      tsp.forEach((day, d) => day.forEach((_, s) => {{
        if (d === di && s === si) items[idx]?.classList.add('active');
        idx++;
      }}));
    }}
  </script>
</body>
</html>"""


@server.list_tools()
async def handle_list_tools() -> list[types.Tool]:
    return [
        types.Tool(
            name="plan_trip",
            description=(
                "Plan a multi-day trip and display an interactive Mapbox map showing "
                "the optimized route for each day with site details."
            ),
            inputSchema={
                "type": "object",
                "properties": {
                    "location": {
                        "type": "string",
                        "description": "City, region, or department (e.g. 'Paris', 'Bretagne')",
                    },
                    "days": {
                        "type": "integer",
                        "description": "Number of trip days (1–7)",
                        "minimum": 1,
                        "maximum": 7,
                    },
                    "types": {
                        "type": "string",
                        "description": "Site type to visit (e.g. 'museum', 'restaurant', 'monument')",
                    },
                },
                "required": ["location", "days", "types"],
            },
        )
    ]


@server.call_tool()
async def handle_call_tool(
    name: str, arguments: dict | None
) -> list[types.TextContent | types.ImageContent | types.EmbeddedResource]:
    if name != "plan_trip":
        raise ValueError(f"Unknown tool: {name}")

    args = arguments or {}
    location = args.get("location", "")
    days = args.get("days", 1)
    types_ = args.get("types", "")

    async with httpx.AsyncClient(timeout=30.0) as client:
        resp = await client.post(
            f"{BACKEND_URL}/api/v1/tsp",
            json={"location": location, "days": days, "types": types_},
        )

    if resp.status_code != 200:
        return [types.TextContent(type="text", text=f"Backend error {resp.status_code}: {resp.text}")]

    data = resp.json()
    tsp: list[list] = data.get("TSP", [])
    total_sites = sum(len(d) for d in tsp)

    if total_sites == 0:
        return [types.TextContent(
            type="text",
            text=f"No sites found for location '{location}' with type '{types_}'.",
        )]

    html = _build_html(data)

    return [
        types.EmbeddedResource(
            type="resource",
            resource=types.TextResourceContents(
                uri=f"trip://{location.lower().replace(' ', '-')}/map",
                mimeType="text/html",
                text=html,
            ),
        )
    ]


async def main() -> None:
    async with stdio_server() as (read_stream, write_stream):
        await server.run(
            read_stream,
            write_stream,
            InitializationOptions(
                server_name="t3-trip-planner",
                server_version="0.1.0",
                capabilities=server.get_capabilities(
                    notification_options=NotificationOptions(),
                    experimental_capabilities={},
                ),
            ),
        )


if __name__ == "__main__":
    asyncio.run(main())
