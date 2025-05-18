import pandas as pd
from collections import defaultdict
from ete3 import Tree, TreeStyle, NodeStyle, TextFace, faces, AttrFace, CircleFace
import matplotlib.cm as cm
import matplotlib.colors as mcolors
import time
import random
import os

df = pd.read_csv(r'<csv-file>')

tree_map = defaultdict(set)
for _, row in df.iterrows():
    cid = int(row['cluster_id'])
    pid = row['parent_cluster_id']
    if pd.notna(pid):
        pid = int(pid)
        tree_map[pid].add(cid)

root_nodes = set(df.loc[df['parent_cluster_id'].isna(), 'cluster_id'].astype(int))

cluster_sizes = df['cluster_id'].value_counts().to_dict()

def generate_color_palette(n_colors):
    """Generate a visually distinct color palette"""
    if n_colors <= 20:
        colors = list(mcolors.TABLEAU_COLORS.values())
        if n_colors > len(colors):
            colors.extend(list(mcolors.CSS4_COLORS.values())[:n_colors-len(colors)])
    else:
        colors = []
        for i in range(n_colors):
            h = i / n_colors
            s = 0.7 + random.random() * 0.3
            v = 0.7 + random.random() * 0.3
            rgb = mcolors.hsv_to_rgb([h, s, v])
            hex_color = mcolors.rgb2hex(rgb)
            colors.append(hex_color)
    return colors

colors = generate_color_palette(len(root_nodes))
root_colors = {root: color for root, color in zip(sorted(root_nodes), colors)}

def build_newick(cid, depth=0):
    if depth < 3:
        print(f"{'  '*depth}Parent Cluster: {cid} -> Children: {list(tree_map[cid]) if cid in tree_map else 'None'}")
    
    if cid not in tree_map:
        return f"Cluster_{cid}"
    
    children = list(tree_map[cid])
    return f"({','.join([build_newick(child, depth+1) for child in children])})Cluster_{cid}"

def find_root_for_cluster(cid):
    current = cid
    while True:
        parent = None
        for pid, children in tree_map.items():
            if current in children:
                parent = pid
                break
        if parent is None:
            return current
        current = parent

cluster_to_root = {}
for cid in df['cluster_id'].unique():
    cid = int(cid)
    if pd.notna(cid):
        cluster_to_root[cid] = find_root_for_cluster(cid)

print(f"Number of root nodes: {len(root_nodes)}")
start_time = time.time()
print("Building subtrees...")
subtrees = []
for root_id in sorted(root_nodes):
    newick = build_newick(root_id)
    subtrees.append(newick)
    print(f"Built subtree for root {root_id}")

print(f"Finished building subtrees in {time.time() - start_time:.2f} seconds.")
combined_newick = f"({','.join(subtrees)})GlobalRoot;"

print("Building ete3 Tree...")
tree_start = time.time()
t = Tree(combined_newick, format=1)
print(f"ete3 Tree built in {time.time() - tree_start:.2f} seconds.")

def layout(node):
    if node.name == "GlobalRoot":
        nstyle = NodeStyle()
        nstyle["size"] = 0
        nstyle["hz_line_width"] = 1
        nstyle["vt_line_width"] = 1
        node.set_style(nstyle)
        return
    
    if node.name.startswith("Cluster_"):
        cid = int(node.name.replace("Cluster_", ""))
        root_id = cluster_to_root.get(cid)
        color = root_colors.get(root_id, "#808080")
        
        size = 10
        if cid in cluster_sizes:
            size = min(20, max(5, int(5 + cluster_sizes[cid] / 5)))
        
        nstyle = NodeStyle()
        nstyle["fgcolor"] = color
        nstyle["size"] = size
        nstyle["shape"] = "sphere"
        
        if len(node.children) == 0:
            nstyle["bgcolor"] = color
        else:
            nstyle["bgcolor"] = "white"
        
        node.set_style(nstyle)
        
        if len(node.children) > 0:
            line_face = faces.RectFace(width=50, height=1, fgcolor=color, bgcolor=color)
            faces.add_face_to_node(line_face, node, column=0, position="branch-right")
        
        if node.is_leaf():
            name_face = TextFace(f" {cid}", fsize=8, fgcolor="black")
            faces.add_face_to_node(name_face, node, column=0, position="branch-right")
        elif len(node.children) > 10:
            name_face = TextFace(f"{cid}", fsize=10, fgcolor="black", bold=True)
            faces.add_face_to_node(name_face, node, column=0, position="branch-top")

ts = TreeStyle()
ts.mode = "c" 
ts.show_leaf_name = False
ts.layout_fn = layout
ts.show_branch_length = False
ts.show_branch_support = False
ts.branch_vertical_margin = 20
ts.arc_start = -180
ts.arc_span = 360
ts.min_leaf_separation = 1
ts.root_opening_factor = 0.1
ts.optimal_scale_level = "full"  
ts.scale = 1.5  

ts.optimal_scale_level = "full"  


ts.title.add_face(TextFace("Cluster Dendrogram", fsize=20, bold=True), column=0)

output_dir = os.path.dirname(df.attrs.get('filepath', '.'))
output_file = os.path.join(output_dir, "cluster_dendrogram_modified.png")
print(f"Saving dendrogram to {output_file}")
t.render(output_file, tree_style=ts, w=1200, h=1200, dpi=120)

print("Showing interactive dendrogram plot...")
t.show(tree_style=ts)
