from flask import Flask, request, render_template
import codecs, pymysql
from datetime import date
app = Flask(__name__)

@app.route("/")
def show():
    return render_template("main.html")

@app.route('/result', methods=["POST", "GET"])
def result():
    connection = pymysql.connect(
        host='database-1.cop2pvzm3623.ap-northeast-1.rds.amazonaws.com',
        db='groupwork_db',
        user='test',
        password='111test',
        charset='utf8',
        cursorclass=pymysql.cursors.DictCursor
    )

    try:
        with connection.cursor() as cursor:
            if request.method == "POST":
                name = request.form["name"]
            else:
                name = request.args.get("name")
            sql = """SELECT recipe.title, recipe.image, recipe.recipe_id FROM recipe JOIN material_recipe ON 
            recipe.recipe_id = material_recipe.recipe_id JOIN material ON material_recipe.material_id
            = material.material_id WHERE material.material_name = %s 
            GROUP BY recipe.title"""
            cursor.execute(sql, name)
            cursor.close()
        # Select結果を取り出す
        results = cursor.fetchall()
        return render_template("main.html", name=name, results=results)
    finally:
        connection.close()

@app.route('/recipe')
def recipe():

    recipeID = request.args.get('id')   # クエリの値取得
    
    connection = pymysql.connect(
        host='database-1.cop2pvzm3623.ap-northeast-1.rds.amazonaws.com',
        db='groupwork_db',
        user='test',
        password='111test',
        charset='utf8',
        cursorclass=pymysql.cursors.DictCursor
    )

    try:
        with connection.cursor() as cursor:
            sql = """SELECT distinct material_name, image, indication, cost, url, title FROM material_recipe
            LEFT JOIN material 
            ON material_recipe.material_id = material.material_id
            LEFT JOIN recipe
            ON material_recipe.recipe_id = recipe.recipe_id
            WHERE material_recipe.recipe_id = %s"""
            cursor.execute(sql, recipeID)
            cursor.close()
        # Select結果を取り出す
        results = cursor.fetchall()

        return render_template("recipe.html", results=results)
    finally:
        connection.close()


if __name__ == '__main__':
    app.debug = True
    app.run(host='0.0.0.0')
