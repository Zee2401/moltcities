## 2025-02-21 - Fast Go Image Buffer Rendering & Hex Parsing

**Learning:** `img.Set(x,y,c)` in Go `image.RGBA` and `fmt.Sscanf` for hex color string parsing introduce high per-pixel overhead (~1080ns/hex parse and ~47ms per 1024x1024 canvas pass). Direct `img.Pix` buffer slice manipulation and manual character arithmetic parsing reduce rendering and parsing time by over 10-30x (~32ns/hex parse and ~3.6ms canvas update).
**Action:** Use direct `img.Pix` byte offset operations and fast non-reflection hex parsing in high-throughput image rendering routines.
