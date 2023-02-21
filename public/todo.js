(function($) {
    'use strict';
    $(function() {
        var todoListItem = $('.todo-list');
        var todoListInput = $('.todo-list-input');
        
        $('.todo-list-add-btn').on("click", function(event) {
            event.preventDefault();
        
            var item = $(this).prevAll('.todo-list-input').val();
        
            if (item) {
                $.post("/todos", {name: item}, addItem)
                //todoListItem.append("<li><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
                todoListInput.val("");
            }
        
        });

        $.get('/todos', function(items){
            items.forEach(e => {
                addItem(e)
            });
        });

        var addItem = function(item) {
            if (item.completed == true) {
                todoListItem.append("<li class='completed' id=" + item.id + "><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' checked='checked' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            } else {
                todoListItem.append("<li id=" + item.id + "><div class='form-check'><label class='form-check-label'><input class='checkbox' type='checkbox' />" + item.name + "<i class='input-helper'></i></label></div><i class='remove mdi mdi-close-circle-outline'></i></li>");
            }
            
        }
        
        todoListItem.on('change', '.checkbox', function() {
            let id = $(this).closest("li").attr("id");
            let $self = $(this);
            let complete = true;
            if ($self.attr('checked')) {
                complete = false;
            }
            $.get("complete-todo/" + id + "?complete=" + complete, function(data) {
                if (complete) {
                    $self.removeAttr('checked');
                } else {
                    $self.attr('checked', 'checked');
                }
                
                $self.closest("li").toggleClass('completed');
            })
        });
        
        todoListItem.on('click', '.remove', function() {
            // url : todos/id method : DELETE
            let id = $(this).closest("li").attr("id");
            let $self = $(this);
            $.ajax({
                url: "todos/" + id,
                type: "DELETE",
                success: function(data) {
                    if (data.success) {
                        $self.parent().remove();
                    }
                    
                }
            })
        });
    
    });
    })(jQuery);