import sys
import os
from collections import defaultdict

def parse_coverage(file_path):
    # lines in coverage.out look like:
    # mode: set
    # github.com/user/repo/pkg/file.go:line.col,line.col numstmts count
    if not os.path.exists(file_path):
        print(f"File {file_path} not found.")
        sys.exit(1)

    file_stats = defaultdict(lambda: {"stmts": 0, "covered": 0})
    blocks = {}
    
    with open(file_path, "r") as f:
        for line in f:
            if line.startswith("mode:"):
                continue
            parts = line.strip().split()
            if len(parts) != 3:
                continue
            
            loc, stmts, count = parts
            stmts = int(stmts)
            count = int(count)
            
            if loc not in blocks:
                blocks[loc] = {"stmts": stmts, "count": count}
            else:
                blocks[loc]["count"] += count

    for loc, info in blocks.items():
        file_name = loc.split(":")[0]
        file_stats[file_name]["stmts"] += info["stmts"]
        if info["count"] > 0:
            file_stats[file_name]["covered"] += info["stmts"]

    return file_stats

def generate_markdown(file_stats, output_file):
    folder_stats = defaultdict(lambda: {"stmts": 0, "covered": 0})
    total_stmts = 0
    total_covered = 0

    for f, stats in file_stats.items():
        folder = os.path.dirname(f)
        folder_stats[folder]["stmts"] += stats["stmts"]
        folder_stats[folder]["covered"] += stats["covered"]
        total_stmts += stats["stmts"]
        total_covered += stats["covered"]

    with open(output_file, "w") as f:
        f.write("# Code Coverage Report\n\n")
        
        # Repo-wise
        repo_cov = (total_covered / total_stmts * 100) if total_stmts > 0 else 0
        f.write("## Repo-Wise Coverage\n\n")
        f.write("| Total Statements | Covered Statements | Coverage % |\n")
        f.write("| --- | --- | --- |\n")
        f.write(f"| {total_stmts} | {total_covered} | {repo_cov:.1f}% |\n\n")

        # Folder-wise
        f.write("## Folder-Wise Coverage\n\n")
        f.write("| Folder | Statements | Covered | Coverage % |\n")
        f.write("| --- | --- | --- | --- |\n")
        for folder in sorted(folder_stats.keys()):
            stats = folder_stats[folder]
            cov = (stats["covered"] / stats["stmts"] * 100) if stats["stmts"] > 0 else 0
            f.write(f"| `{folder}` | {stats['stmts']} | {stats['covered']} | {cov:.1f}% |\n")
        f.write("\n")

        # File-wise
        f.write("## File-Wise Coverage\n\n")
        f.write("| File | Statements | Covered | Coverage % |\n")
        f.write("| --- | --- | --- | --- |\n")
        for file_name in sorted(file_stats.keys()):
            stats = file_stats[file_name]
            cov = (stats["covered"] / stats["stmts"] * 100) if stats["stmts"] > 0 else 0
            f.write(f"| `{file_name}` | {stats['stmts']} | {stats['covered']} | {cov:.1f}% |\n")
            
if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python generate_coverage_report.py <coverage.out> <output.md>")
        sys.exit(1)
    
    stats = parse_coverage(sys.argv[1])
    generate_markdown(stats, sys.argv[2])
    print(f"Generated {sys.argv[2]}")
