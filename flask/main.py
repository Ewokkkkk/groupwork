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
            name = ""
            data = ""
            data2 = ""
            if request.method == "POST":
                data = request.form["name"]
            else:
                if request.args.get("category") != None:
                    data = request.args.get("category")
                    sql = """SELECT recipe.title, recipe.image, recipe.recipe_id FROM recipe JOIN category_recipe ON 
                    recipe.recipe_id = category_recipe.recipe_id JOIN category_list ON category_recipe.category_id
                    = category_list.category_id WHERE %s IN(category_list.category_id, category_list.parent_category_id)
                    GROUP BY recipe.title"""
                    cursor.execute(sql, data)
                elif request.args.get("category") == None and request.args.get("ignore") == None:
                    data = request.args.get("name")
                    d_list = data.split(",")
                    name = data
                    sql = """SELECT recipe.title, recipe.image, recipe.recipe_id FROM recipe 
                    JOIN material_recipe 
                    ON recipe.recipe_id = material_recipe.recipe_id 
                    JOIN material 
                    ON material_recipe.material_id = material.material_id 
                    WHERE recipe.recipe_id IN 
                    ( SELECT recipe_id
                    FROM material_recipe
                    WHERE material_id IN
					( SELECT material.material_id
					FROM material
					WHERE material_name LIKE %s
                    ))"""
                    
                    for i in range(len(d_list)-1):
                        sql += """AND recipe.recipe_id IN
                    ( SELECT recipe_id
                    FROM material_recipe
                    WHERE material_id IN
					( SELECT material.material_id
					FROM material
					WHERE material_name LIKE %s
                    ))"""

                    sql += """GROUP BY recipe.title"""

                    for i in range(len(d_list)):
                        d_list[i] = "%" + d_list[i] + "%"

                    cursor.execute(sql, d_list)

                else:
                    data = request.args.get("name")
                    d_list = data.split(",")
                    data2 = request.args.get("ignore")
                    d2_list = data2.split(",")
                    name = data
                    # sql = """SELECT recipe.title, recipe.image, recipe.recipe_id FROM recipe 
                    # JOIN material_recipe 
                    # ON recipe.recipe_id = material_recipe.recipe_id 
                    # JOIN material 
                    # ON material_recipe.material_id = material.material_id 
                    # WHERE recipe.recipe_id NOT IN 
                    # ( SELECT recipe_id
                    # FROM material_recipe
                    # WHERE material_id IN
					# 	( SELECT material.material_id
					# 	FROM material
					# 	WHERE material_name LIKE %s)
                    # ) AND material.material_name LIKE %s 
                    # GROUP BY recipe.title"""
                    # sql = """SELECT recipe.title, recipe.image, recipe.recipe_id FROM recipe 
                    # JOIN material_recipe 
                    # ON recipe.recipe_id = material_recipe.recipe_id 
                    # JOIN material 
                    # ON material_recipe.material_id = material.material_id 
                    # WHERE recipe.recipe_id NOT IN 
                    # ( SELECT recipe_id
                    # FROM material_recipe
                    # WHERE material_id IN
					# ( SELECT material.material_id
					# FROM material
					# WHERE material_name LIKE %s """
                    # for i in range(len(d2_list)-1):
                    #     sql += """ OR material.material_name LIKE %s """

                    # sql += """ )
                    # ) AND material.material_name LIKE %s """

                    # for i in range(len(d_list)-1):
                    #     sql += """ OR material.material_name LIKE %s """

                    sql = """SELECT recipe.title, recipe.image, recipe.recipe_id FROM recipe 
                    JOIN material_recipe 
                    ON recipe.recipe_id = material_recipe.recipe_id 
                    JOIN material 
                    ON material_recipe.material_id = material.material_id 
                    WHERE recipe.recipe_id NOT IN 
                    ( SELECT recipe_id
                    FROM material_recipe
                    WHERE material_id IN
					( SELECT material.material_id
					FROM material
					WHERE material_name LIKE %s
                    ))"""
                    for i in range(len(d2_list)-1):
                        sql += """AND recipe.recipe_id NOT IN
                    ( SELECT recipe_id
                    FROM material_recipe
                    WHERE material_id IN
					( SELECT material.material_id
					FROM material
					WHERE material_name LIKE %s
                    ))"""

                    for i in range(len(d_list)):
                        sql += """AND recipe.recipe_id IN
                    ( SELECT recipe_id
                    FROM material_recipe
                    WHERE material_id IN
					( SELECT material.material_id
					FROM material
					WHERE material_name LIKE %s
                    ))"""

                    for i in range(len(d_list)):
                        d_list[i] = "%" + d_list[i] + "%"
                    for i in range(len(d2_list)):
                        d2_list[i] = "%" + d2_list[i] + "%"

                    sql += """GROUP BY recipe.title"""


                    print(sql)
                    cursor.execute(sql, d2_list + d_list)

            # cursor.execute(sql, data)
        # Select結果を取り出す 
            results = cursor.fetchall()
            if name == "":
                sql = """SELECT category_name FROM category_list WHERE category_list.category_id = %s"""
                cursor.execute(sql, data)
                data = cursor.fetchone()
                data = data["category_name"]
                print(data)
            cursor.close()
        if data2 != "":
            return render_template("result.html", name=data, results=results, ng=data2)
        else :
            return render_template("result.html", name=data, results=results)

    finally:
        connection.close()

@app.route('/recipe')
def recipe():

    recipeID = request.args.get('id')   # クエリの値の取得
    
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
        
@app.route('/list')
def m_list():
    return render_template("list.html")

if __name__ == '__main__':
    app.debug = True
    app.run(host='0.0.0.0')